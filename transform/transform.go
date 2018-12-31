package transform

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
)

const PIPE_SIZE = 65536

type Transformer interface {
	Eval(code string) (result interface{}, error)
	Eval(name string, val interface{}) error
}

// backup

// func (c *BaseControl) transformValue() interface{} {
// 	if c.Transform != "" {
// 		if c.transform == nil {
// 			c.transform = &transform.Transformer{
// 				Program: "./transform.js",
// 			}
// 			c.transform.Init()
// 		}

// 		c.transform.Set("v", c.value)
// 		newVal, err := c.transform.Eval(c.Transform)

// 		if err != nil {
// 			c.setError(err)
// 		}

// 		return newVal
// 	} else {
// 		return c.value
// 	}
// }
