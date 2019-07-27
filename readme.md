# PicBed

Pic bed in golang. This golang package consist all open picture bed uploading API enable you upload image to those picture bed.



Projects using PicBed:

- dingo: http://jinfagang/dingo





## Usage

For using PicBed. You can call those API for now:

- `Ali`:

  ```go
  imgBytes, err := ioutil.ReadFile(imgPath)
  if(err != nil) {
  fmt.Println(err)
  }
  imgParam := picbed.ImageParam{
  Name:    filepath.Base(imgPath),
  Type:    "jpg",
  Content: &imgBytes,
  }
  // using Ali as default
  client := picbed.Ali{FileLimit: nil, MaxSize: 5024}
  res, _ := client.Upload(&imgParam)
  fmt.Println(res.Url)
  ```

- `Jd`:

  ```
  to be done
  ```

- `Gitee`:

  ```
  to be done
  ```

  

## Copyright

All right belongs to Jin Fagang, codes released under Apache License.