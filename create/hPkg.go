package create

import (
	"github.com/cgalvisleon/elvis/file"
	"github.com/cgalvisleon/elvis/strs"
)

func MakePkg(name, schema string) error {
	path, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "event.go", modelEvent, name)
	if err != nil {
		return err
	}

	modelo := strs.Titlecase(name)
	_, err = file.MakeFile(path, "model.go", modelModel, name, modelo)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "msg.go", modelMsg, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "controller.go", modelController, name)
	if err != nil {
		return err
	}

	if len(schema) > 0 {
		err = MakeDbHandler(name, name, schema)
		if err != nil {
			return err
		}

		title := strs.Titlecase(name)
		_, err = file.MakeFile("router.go", modelDbRouter, name, title)
		if err != nil {
			return err
		}
	} else {
		err = MakeHandler(name, name)
		if err != nil {
			return err
		}

		title := strs.Titlecase(name)
		_, err = file.MakeFile("router.go", modelRouter, name, title)
		if err != nil {
			return err
		}
	}

	return nil
}

func MakeDbHandler(name, modelo, schema string) error {
	path, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	schemaVar := strs.Append("schema", strs.Titlecase(schema), "")

	_, err = file.MakeFile(path, "schema.go", modelSchema, name, schemaVar, schema)
	if err != nil {
		return err
	}

	path, err = file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = strs.Titlecase(modelo)
	fileName := strs.Format(`h%s.go`, modelo)
	_, err = file.MakeFile(path, fileName, modelDbHandler, name, modelo, schemaVar, strs.Uppcase(modelo), strs.Lowcase(modelo))
	if err != nil {
		return err
	}

	return nil
}

func MakeHandler(name, modelo string) error {
	path, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	modelo = strs.Titlecase(modelo)
	fileName := strs.Format(`h%s.go`, modelo)
	_, err = file.MakeFile(path, fileName, modelHandler, name, modelo, strs.Lowcase(modelo))
	if err != nil {
		return err
	}

	return nil
}

func MakeRpc(name string) error {
	path, err := file.MakeFolder("pkg", name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "hRpc.go", modelhRpc, name)
	if err != nil {
		return err
	}

	return nil
}
