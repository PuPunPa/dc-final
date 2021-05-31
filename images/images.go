package images
import(
    "gocv.io/x/gocv"
	"github.com/jeasonstudio/GaussianBlur"
)
func Image2GrayScale(path string){
    mat:=gocv.IMRead(path,gocv.IMReadGrayScale)
    if!mat.Empty(){
        gocv.IMWrite(path+"GrayScales.png",mat)
    }
}
func Image2Blur(path string)  {
    GaussianBlur.GBlurInit(path,path+"blurred.png",5,5.0)
}