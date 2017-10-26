package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type TypeEnv struct {
	DB_DRIVER string `json:"DB_DRIVER"`
	DB_ARGS   string `json:"DB_ARGS"`
	SV_URL    string `json:"SV_URL"`
	DEBUG     string `json:"DEBUG"`
}

var Env = TypeEnv{}

func ReadEnv(path string) error {
	data, err := readFile(path)
	if err != nil {
		return err
	}
	fmt.Printf("[CONF]\n%s\n", data)
	err = json.Unmarshal(data, &Env)
	return err
}

func readFile(path string) ([]byte, error) {
	fin, er := os.Open(path)
	if er != nil {
		fmt.Printf("failure to open the file : %s\n", path)
		return nil, er
	}
	defer fin.Close()
	buf, er := ioutil.ReadAll(fin)
	if er != nil {
		fmt.Print("failure to read the file : %s\n", path)
		return nil, er
	}
	return buf, nil
}
