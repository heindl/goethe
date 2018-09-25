// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cobra

import (
	"go/ast"
	"path"
	"sync"

	"github.com/heindl/goethe/utilities"
	"golang.org/x/sync/errgroup"
)

type cobraRootCommands map[string][]string

// TODO: This is a implementation should be refactored with the ssa package for simplicity and correctness ...
// Current implementation will miss function reassignments and recursive functions.
// But the tooling around ssa is still in active development.

func search(info *utilities.ModuleInfo) (cobraRootCommands, error) {

	// Note that Packages will only call ParseDir once across the runtime.
	pkgs, err := info.Packages()
	if err != nil {
		return nil, err
	}

	y := cobraRootCommands{}
	locker := sync.Mutex{}

	eg := errgroup.Group{}
	for _, _pkg := range pkgs {
		pkg := _pkg
		eg.Go(func() error {
			vars := cobraCmdVars{vars: make(map[string]struct{})}
			ast.Walk(vars, pkg.AstPkg)
			execs := execCalls{
				importPath: pkg.ImportPath,
			}
			ast.Walk(&execs, pkg.AstPkg)
			for _, ex := range execs.execs {
				if _, ok := vars.vars[path.Base(ex.cmdVarName)]; ok {
					locker.Lock()
					if _, ok := y[ex.cmdVarName]; !ok {
						y[ex.cmdVarName] = []string{}
					}
					y[ex.cmdVarName] = append(y[ex.cmdVarName], path.Dir(ex.enclFunc))
					locker.Unlock()
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return y, nil
}

type execCalls struct {
	enclFunc   string
	execs      []execCall
	importPath string
}

type execCall struct {
	cmdVarName string
	enclFunc   string
}

func (Ω *execCalls) Visit(nd ast.Node) ast.Visitor {

	if nd == nil {
		return Ω
	}

	switch v := nd.(type) {
	case *ast.File:
		Ω.enclFunc = ""
		return Ω
	case *ast.FuncDecl:
		Ω.enclFunc = v.Name.String()
		return Ω
	case *ast.SelectorExpr:
		// Pass through
	default:
		return Ω
	}

	selectorExpr := nd.(*ast.SelectorExpr)

	if selectorExpr.Sel.Name != "Execute" {
		return Ω
	}

	id, isIdent := selectorExpr.X.(*ast.Ident)
	if !isIdent {
		return Ω
	}

	Ω.execs = append(Ω.execs, execCall{
		cmdVarName: path.Join(Ω.importPath, id.String()),
		enclFunc:   path.Join(Ω.importPath, Ω.enclFunc),
	})

	return Ω
}

type cobraCmdVars struct {
	vars          map[string]struct{}
	lastValueName string
}

func (Ω cobraCmdVars) Visit(nd ast.Node) ast.Visitor {
	valueSpec, isValueSpec := nd.(*ast.ValueSpec)
	if isValueSpec {
		Ω.lastValueName = valueSpec.Names[0].String()
		return Ω
	}

	selector, isSelector := nd.(*ast.SelectorExpr)
	if !isSelector || selector.Sel.String() != "Command" {
		return Ω
	}
	id, isIdent := selector.X.(*ast.Ident)
	if !isIdent || id.String() != "cobra" {
		return Ω
	}
	Ω.vars[Ω.lastValueName] = struct{}{}
	return Ω
}
