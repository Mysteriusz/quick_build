package clang_vc

import(
	"fmt"

	"qb/misc"
	"qb/qbio"
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/cfg"
)

/*
	Verify output working set object
	based on the clang object storage
*/
func OutObjectVerify(_obj qb.Object, _vc_state *vc.FileState) (removed bool, no_change bool){
	obj_file := _obj.Data.(qb.FileObject).File

	/*
		Resolve extra objects for the file
	*/
	_,dep_file := qb.GetObjectExtra[qbio.File](&_obj, clang.OUT_DEP)
	_,src_file := qb.GetObjectExtra[qbio.File](&_obj, clang.OUT_SRC)

	/*
		Validate the Source file
	*/

	v := !src_file.IsValid() 
	h := !src_file.InvalidateHash()
	if v || h{
		fmt.Printf("Source file changed:\n %s\n", src_file.FullPath)
		return v, false // If invalid then source is removed
	}

	/*
		Validate the Dependency file
	*/

	v = !dep_file.IsValid() 
	h = !dep_file.InvalidateHash()
	if v || h{
		fmt.Printf("Dependency file changed:\n %s\n", dep_file.FullPath)
		return false, false
	}

	/*
		Validate the Object file
	*/

	v = !obj_file.IsValid() 
	h = !obj_file.InvalidateHash()
	if v || h{
		fmt.Printf("Object file changed:\n %s\n", obj_file.FullPath)
		return false, false
	}

	/*
		Parse the dependency file
	*/

	res, file := clang.ParseD(dep_file)
	if !res{
		fmt.Printf("Unable to parse the dependency file:\n %s\n", dep_file.FullPath)
		return false, false
	}

	/*
		Check if object file is the same as the dependency file one
	*/

	if !obj_file.Compare(file.Deps[0]){
		fmt.Printf("Dependency source files do not match:\n %s\n %s\n", obj_file.FullPath, file.Deps[0].FullPath)
		return false, false
	}

	/*
		Check if source file is the same as the dependency file one
	*/

	if !src_file.Compare(file.Deps[1]){
		fmt.Printf("Dependency source files do not match:\n %s\n %s\n", src_file.FullPath, file.Deps[1].FullPath)
		return false, false
	}

	/*
		Intersect headers from the dependency file and gathered by the 'vc.FileState'
		If any are matched that means the file has to be re-compiled
	*/
	shared_diffs := misc.Intersect(file.Deps[2:].AllPaths(), _vc_state.DiffHeaders.Modified.AllPaths())
	if len(shared_diffs) != 0{
		fmt.Println("Header files do not match")
		return false, false
	}

	return false, true
}

