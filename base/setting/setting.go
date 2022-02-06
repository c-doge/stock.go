package setting

import (
    yaml "gopkg.in/yaml.v3"
    "io/ioutil"
    "fmt"
)

type Setting struct {
    Web struct {
        Addr    string  `yaml:"addr"`
    }
    Log struct {
        Path    string  `yaml:"path"`
        Level   string  `yaml:"level"`
    }
    LevelDB struct {
        Path    string  `yaml:"path"`
    }
}

var settings *Setting = new(Setting)

func Parse(filepath string) error {

    if settings == nil {
        settings = new(Setting)
    }

    yamlFile, err := ioutil.ReadFile(filepath)

    if err != nil {
        return err;
    }

    err = yaml.Unmarshal(yamlFile, settings)
    if err != nil {
        return err
    }

    return nil
}

func Get() *Setting {
    return settings;
}

func (*Setting) String() string {
    b, err := yaml.Marshal(settings);
    if err != nil {
        return fmt.Sprintf("yaml Marshal Fail, err: %v", err);    
    }
    return string(b) 
}