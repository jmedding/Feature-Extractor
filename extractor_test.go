package main 

import (
  "testing"
  "image"
  "os"
  "path"
  "strings"
  "path/filepath"
  //"fmt"
)

func clean(fpath string, f os.FileInfo, err error) error{
  //delete the directory if contains 'testing'
  //err will not be nil as deleting directories will make open
  //their sub-files/directories impossible. 
  //No big deal, just don't return the error

  match := strings.Contains(fpath, "testing") 
  //fmt.Println(err, match, fpath)
  if match && err == nil{
    os.RemoveAll(fpath)
  }
  return nil
}

func Test_visit( t * testing.T){
  //using 'testing' folder, located in extractor src dir for images
  //testing location is not ideal, but don't know how to move it out of 'test_iamges'
  //without breaking extractor.go - badd design...
  image_path := path.Join(base_dir(), "test_images"  ,"testing/images/group_1/1.jpeg")  //contains group_1{group_1_1}, groupi_2
  
  fileInfo, err := os.Stat(image_path)
  if err != nil {
    t.Error("LStat could not return os.FileInfo. error:", err)
  } else{
    err2 := visit(image_path, fileInfo, nil)
    if err2 != nil {
      t.Error("visit() returned and error:", err2)
    } else {
      //check that features were extracted and written in proper location
      output_path := strings.Replace(image_path, "images", "training", 1)
      found, err3 := os.Open(image_path)
      defer found.Close()
      if err3 != nil {
        t.Error ("expected to find file at", output_path)
      } else {
        t.Log(".")
        //clean up created directories
        //deleate any directory called 'testing' in  'training' dir
        filepath.Walk(path.Join(base_dir(), "training"), clean )

      }
    }
  }
}

func Test_get_face(t *testing.T){
  //should return a feature, with (X,Y) as the center of the face in the image
  r, _ := get_response("__testing__")
  ist := r.get_face()
  soll := Feature{97, 133, 133, 133}
  if soll != ist {
    t.Error("get_face() returned", ist, "but should be", soll)
  } else {
    t.Log(".")
  }
}

func Test_get_feature_abs(t *testing.T){
  //should return the a rectangle for the feature in the image coords
  r, _ := get_response("__testing__")
  eyeL := r.get_feature_abs(r.Eye_left)
  soll := image.Rect(44,104, 94, 130)
  if soll != eyeL {
    t.Error("get_feature_abs() returned", eyeL, "but should be", soll)
  } else {
    t.Log(".")
  }
}

func Test_Rect(t *testing.T) {

  //feate has X,Y at center of feature
  f := Feature { 
    X:1,
    Y:2,
    Width: 2,
    Height: 4,
  }

  rectangle := f.Rect()
  min := rectangle.Min
  if min != image.Pt(0, 0) {
    t.Error("Feature.Rect() returned wrong Min:", min)
  } 
  max := rectangle.Max
  if  max != image.Pt(2,4){
    t.Error("Feature.Rect() returned wrong Max:" , max)
  } else{
    t.Log(".")
  }
}