package main

import (
  "fmt"
  "path/filepath"
  "os"
  "path"
  "flag"
  "strings"
  "encoding/json"
  "image"
  "image/draw"
  "image/jpeg"
  _"image/png"
)
/*
This program should work through a list of pictures,
In each picture, it should identify a face
In that face it should also locate the 
  -left and right eyes
  -mouth
  -top of the head/hair
It can do this by making a call to a detection program 
that returns the coordinates of the above items.

It should then extract portions of the images that correspond
to the detected features and save these in a useful structure

*/
type Feature struct{
  X, Y, Width, Height int
}

func base_dir() string { return "/home/jon/face_js"}

func dest_dir() string { return "training"}

func (f Feature) Rect() (image.Rectangle){
  return image.Rect(f.X - f.Width/2, f.Y-f.Height/2, f.X + f.Width/2, f.Y + f.Height/2)
}

type Response struct{
  Status string
  Messages string
  Rotation int
  X, Y, Width, Height int
  Eye_right, Eye_left, Mouth  Feature
}

func (r Response) get_face() Feature{
  return Feature{r.X, r.Y, r.Width, r.Height}
}

func (r Response) get_feature_abs(f Feature) (image.Rectangle) {
  //will return the feature rectangle in the image coord
  //do not use with get_face() value - it is already in the image coords
  return f.Rect().Add(image.Pt(r.X - r.Width/2, r.Y - r.Height/2))
}

func dummy_response() string {
  return "{ \"status\": \"success\",\"messages\": \"\", \"rotation\": 0, \"x\": 97, \"y\": 133, \"height\": 133, \"width\": 133, \"eye_left\" : { \"x\": 38, \"y\": 50, \"height\": 26, \"width\": 50}, \"eye_right\" : { \"x\": 93, \"y\": 50, \"height\": 26, \"width\": 40}, \"mouth\" : { \"x\": 65, \"y\": 114, \"height\": 26, \"width\": 44}} "
}

func get_response(path string) (Response, error) {
  var r Response
  //testing case
  raw_response := dummy_response()
  err := json.Unmarshal([]byte(raw_response), &r)
  if err != nil {
    fmt.Println(err)
  }
  return r, err
}

func get_image(path string) (image.Image, error){
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil {
    return nil, err
  }

  // fmt.Printf("Image type: %T", img)
  return img, err
}

func CreateRGBAfromImage(img image.Image) *image.RGBA {
  rect := img.Bounds()
  rgba_image := image.NewRGBA(rect)
  draw.Draw(rgba_image, rect, img, rect.Min, draw.Src)
  return rgba_image
}

func write_image(fpath, sub_dir string, img image.Image) error{
  dir, name := path.Split(fpath)
  dir_parts := strings.Split(dir, "/")
  base_dir_parts := strings.Split(base_dir(), "/")
  base_dir_parts_len := len(base_dir_parts)
  // fmt.Println(base_dir_parts, len(base_dir_parts))
  // fmt.Println(dir_parts[base_dir_parts_len + 1:], len(dir_parts))
  source_dir := strings.Join(dir_parts[ base_dir_parts_len + 1:], "/")
  // fmt.Println(source_dir)
  write_dir := path.Join(base_dir(), dest_dir(), sub_dir, source_dir)
  os.MkdirAll(write_dir, 0777)
  write_path := path.Join(write_dir, name)
  // fmt.Printf("%T %s\n", img, img.Bounds())
  os.MkdirAll(write_dir, 0777)
  dest_file, err := os.Create(write_path)
  defer dest_file.Close()
  if err != nil {
    fmt.Println(err)
    return err
  }
  // fmt.Println(write_dir)
  jpeg.Encode(dest_file, img, &jpeg.Options{jpeg.DefaultQuality})
  return nil
}

func is_image(fpath, fname string) {
  // fmt.Printf("Processing: %s", fname)
  img, err := get_image(fpath)
  if err != nil {
    fmt.Println(err)
  }
  // fmt.Println(img.Bounds())
  r, err := get_response(fpath)
  if err != nil {
    fmt.Println(err)
  }
  // jpeg and png can decode to differnt color models
  // so fource everything to RGBA
  rgba := CreateRGBAfromImage(img)
  face := rgba.SubImage(r.get_face().Rect())
  write_image(fpath, "face", face)
  eyeL := rgba.SubImage(r.get_feature_abs(r.Eye_left))
  write_image(fpath, "eye_left", eyeL)
  eyeR := rgba.SubImage(r.get_feature_abs( r.Eye_right))
  write_image(fpath, "eye_right", eyeR)
  mouth := rgba.SubImage(r.get_feature_abs(r.Mouth))
  write_image(fpath, "mouth", mouth)
}

func visit(fpath string, f os.FileInfo, err error) error {
  if f.IsDir() && f.Name() == "big" {
    fmt.Printf("Visited: %s\n", fpath)
  }

  if strings.HasSuffix(fpath, ".jpg") ||
  strings.HasSuffix(fpath, ".jpeg") ||
  strings.HasSuffix(fpath, ".png") {
    is_image(fpath, f.Name())
  }
  return nil
}

func main() {
    flag.Parse()
    fmt.Printf("I'm gonna get some faces!\n")
    root := path.Join (base_dir(), "test_images", flag.Arg(0))
    fmt.Println(root)
    err := filepath.Walk(root, visit)
    fmt.Printf("filepath.Walk() returned %v\n", err)
}