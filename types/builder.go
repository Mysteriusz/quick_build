package types

import(
	"path/filepath"
	"io/fs"
	"slices"
	"errors"

	"qb/qberr"
)

type EntryBuilder struct{
	EntryKey 	string
	doc 		*Doc
	entry 		*Entry
	ver 		BuildVersion
	base_dir 	string // Resolved
	out_dir 	string // Resolved
	compiler_cmd 	*Command
	linker_cmd	*Command
	library_cmd	*Command
	compiled 	[]string
}

func (_builder *EntryBuilder) GetDoc() *Doc{
	if _builder == nil || _builder.doc == nil{
		panic("Entry builder is corrupted!")
	}
	return _builder.doc
}
func (_builder *EntryBuilder) GetEntry() *Entry{
	if _builder == nil || _builder.doc == nil{
		panic("Entry builder is corrupted!")
	}
	return _builder.GetDoc().FromKey(_builder.EntryKey)
}

func CreateEntryBuilder(_doc *Doc, _key string) (*EntryBuilder, bool){
	if _doc == nil{
		return nil, false
	}

	var res bool
	var entry *Entry = _doc.FromKey(_key)
	var builder *EntryBuilder = &EntryBuilder{}

	builder.EntryKey = _key

	builder.doc = _doc
	builder.entry = entry

	/*
		Resolve and store directories
	*/
	builder.base_dir, res = _doc.PathFromKey(_key, entry.BaseDirectory)
	if !res{
		return nil, false
	}

	builder.out_dir, res = _doc.PathFromKey(_key, entry.BuildDirectory)
	if !res{
		return nil, false
	}

	/*
		(MANDATORY)
		Resolve and store default compiler command
	*/
	builder.compiler_cmd, res = 
		_doc.Compiler.CompilerPrefixGroup.DefaultCmd(builder)
	if !res{
		return nil, false
	}

	/*
		(OPTIONAL)
		Resolve and store default linker command
	*/
	if _doc.Compiler.LinkerPrefixGroup.NAME != ""{
		builder.linker_cmd, res =
			_doc.Compiler.LinkerPrefixGroup.DefaultCmd(builder)
		if !res{
			return nil, false
		}
	}

	/*
		(OPTIONAL)
		Resolve and store default library creator command
	*/
	if _doc.Compiler.LibraryPrefixGroup.NAME != ""{
		builder.library_cmd, res = 
			_doc.Compiler.LibraryPrefixGroup.DefaultCmd(builder)
		if !res{
			return nil, false
		}
	}

	/*
		Load version control document if entry requests it
	*/
	if entry.VersionControl{
		res = build_version_load(builder)
		if !res{
			return nil, false
		}
	}

	builder.compiled = make([]string, 100)

	return builder, true
}

func (_builder *EntryBuilder) Compile() bool{
	if _builder == nil{
		return false
	}

	var res bool

	doc := _builder.GetDoc()
	entry := _builder.GetEntry()

	// Walk over all entry files and compile them
	filepath.Walk(filepath.Join(entry.BaseDirectory, "src"), func(path string, info fs.FileInfo, err error) error{
		if err != nil{
			qberr.ERR("An error occured.", err)
			return err
		}
    		if info.IsDir(){
			return nil
		}

		var out_path string
		var mode CommandMode
		/*
			Check if file is neither a source file
		*/
		if slices.Contains(doc.Compiler.SourceExtensions, filepath.Ext(path)){
			out_path = filepath.Join(_builder.out_dir, info.Name() + doc.Compiler.SourceOutExtension)
			mode = Source
		/*
			Check if file is neither a source file
		*/
		}else if slices.Contains(doc.Compiler.HeaderExtensions, filepath.Ext(path)){
			out_path = filepath.Join(_builder.out_dir + "headers", info.Name()); 
			mode = Header
		/*
			Unknown file extension (Ignore)
		*/
		}else{
			return nil
		}

		if entry.VersionControl{
			_, res = build_version_check_and_update(_builder, path)
			if !res{
				return errors.New("")
			}
		}

		_builder.compiler_cmd.Mode = mode 
		_builder.compiler_cmd.Input = path
		_builder.compiler_cmd.Output = out_path

		res = _builder.doc.Compiler.CompilerPrefixGroup.FromCmd(_builder.compiler_cmd)
		if !res{
			return errors.New("")
		}

		return nil
	})

	return true
}

func (_builder *EntryBuilder) Process(){
}

func (_builder *EntryBuilder) Finish() bool{
	if _builder == nil{
		return false
	}

	var res bool
	entry := _builder.GetEntry()

	if entry.VersionControl{
		res = build_version_save(_builder)
		if !res{
			return false
		}
	}
	return true
}

