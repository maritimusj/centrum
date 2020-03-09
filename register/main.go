package register

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/maritimusj/centrum/util"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/net"
)

var (
	nonce             = "20200102"
	fingerprintsCache = ""
)

//硬件指纹
func fingerprints(nonce string) []byte {
	var data bytes.Buffer

	v, _ := cpu.Info()
	for _, x := range v {
		if x.PhysicalID != "" {
			data.WriteString(x.PhysicalID)
		}
	}

	if data.Len() < 1 {
		v, _ := net.Interfaces()
		for _, x := range v {
			if !strings.Contains(x.Name, "Loopback") && !strings.Contains(x.Name, "isatap") && x.HardwareAddr != "" {
				data.WriteString(x.HardwareAddr)
			}
		}
	}

	if data.Len() > 0 {
		data.WriteString(nonce)
		hash := sha256.Sum256(data.Bytes())
		return hash[:]
	}

	return nil
}

//获取本机特征码
func Fingerprints() string {
	if fingerprintsCache != "" {
		return fingerprintsCache
	}

	fingerprintsCache = hex.EncodeToString(fingerprints(nonce))
	return fingerprintsCache
}

//计算一个注册码
/*
	第一步：生成一个4位随机密码
	第二步：使用上面的密码做为hmac算法密码加密（owner + fingerprints)
	第三步：取第二步结果的首尾各4位字符,
	第四步：密码+第三步的两个字符串合成注册码
*/
func Code(owner, hardwareFingerprints string) string {
	randStr := strings.ToLower(util.RandStr(6, util.RandAll))
	hash := hmac.New(sha1.New, []byte(randStr))
	hashStr := hex.EncodeToString(hash.Sum([]byte(owner + hardwareFingerprints)))

	return fmt.Sprintf("%s-%s-%s", randStr, hashStr[:6], hashStr[len(hashStr)-6:])
}

//检验注册码
func Verify(owner, code string) bool {
	if code != "" {
		codes := strings.Split(code, "-")
		if len(codes) == 3 {
			hash := hmac.New(sha1.New, []byte(codes[0]))
			x := hex.EncodeToString(hash.Sum([]byte(owner + Fingerprints())))
			return x[:6] == codes[1] && x[(len(x)-6):] == codes[2]
		}
	}

	return false
}
