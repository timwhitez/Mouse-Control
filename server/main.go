package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"math/big"
	r1 "math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	CmdStr     string
	ReShellStr string
	Host       = "0.0.0.0" //服务器ip
	Port       = "18088"
)

var A0 string = ""

// ReShell
func getmsg(c *gin.Context) {
	for A0 == "" {
	}
	c.String(200, A0)
	A0 = ""
	//c.String(404,"")
}

func GetRandomString(n int) string { //随机字符串
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_@!#$%^"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[r1.Intn(len(bytes))])
	}
	return string(result)
}

func CreatPem(host string, pemname string, keyname string) {
	max := new(big.Int).Lsh(big.NewInt(1), 256)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{GetRandomString(6)},
		OrganizationalUnit: []string{GetRandomString(5)},
		CommonName:         GetRandomString(10),
	}

	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(36500 * time.Hour),
		//KeyUsage:     x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP(host)},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &pk.PublicKey, pk)

	certOut, _ := os.Create(pemname)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	certOut.Close()

	keyOut, _ := os.Create(keyname)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}

func CheckFileExist(fileName string) bool { //检查配置文件是否存在
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     Host + ":" + Port,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}

func control(func0 string, speed int) {
	A0 = func0 + "|" + strconv.Itoa(speed)
}

var input1 *widget.Entry

// 输入
func Inp() *fyne.Container {

	//文件名输入框
	input1 = widget.NewEntry()
	input1.SetPlaceHolder("Enter text...")
	input1.Wrapping = fyne.TextWrapOff

	//将上述排版
	c := container.NewHBox(widget.NewLabel("speed: "), input1)

	return c
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	//gin.DisableConsoleColor()
	fmt.Println(Host+":"+Port)
	go func() {
		router := gin.Default()
		gin.SetMode(gin.ReleaseMode)
		router.Use(TlsHandler())
		router.GET("/args", getmsg)

		if CheckFileExist("./tls.pem") && CheckFileExist("./tls.key") {
			err := router.RunTLS(Host+":"+Port, "./tls.pem", "./tls.key")
			if err != nil {
				fmt.Println(err)
			}
		} else {
			CreatPem(Host, "./tls.pem", "./tls.key")
			err := router.RunTLS(Host+":"+Port, "./tls.pem", "./tls.key")
			if err != nil {
				fmt.Println(err)
			}
		}

	}()

	a := app.New()
	w := a.NewWindow("Rob")
	//hello := widget.NewLabel("Hello Fyne!")

	w.SetContent(container.NewHBox(layout.NewSpacer(), widget.NewLabel("      "), container.NewVBox(
		//hello,
		container.NewHBox(layout.NewSpacer(),
			widget.NewButton("LC", func() {
				control("LC", 0)
			}),
			widget.NewButton("DouC", func() {
				control("DouC", 0)
			}),
			widget.NewButton("RC", func() {
				control("RC", 0)
			}), layout.NewSpacer()),

		Inp(),
		container.NewHBox(layout.NewSpacer(), widget.NewButton("  up  ", func() {
			i, _ := strconv.Atoi(input1.Text)
			control("up", i)
		}),
			layout.NewSpacer()),
		container.NewHBox(
			widget.NewButton(" left ", func() {
				i, _ := strconv.Atoi(input1.Text)
				control("left", i)
			}),
			widget.NewLabel("            "),
			widget.NewButton("right", func() {
				i, _ := strconv.Atoi(input1.Text)
				control("right", i)
			}),
		),
		container.NewHBox(layout.NewSpacer(),
			widget.NewButton(" down ", func() {
				i, _ := strconv.Atoi(input1.Text)
				control("down", i)
			}),
			layout.NewSpacer()),
	), widget.NewLabel("      "), layout.NewSpacer()))

	w.ShowAndRun()

}
