package main

import "os"

type Envs struct{}

func (e *Envs) StringVar(dst *string, def string, varname string) {
	if *dst != def {
		return
	}
	v := os.Getenv(varname)
	if v != "" {
		*dst = v
	}
}
