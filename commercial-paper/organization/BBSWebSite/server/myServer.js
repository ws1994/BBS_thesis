let express = require("express")
let http = require("http")

let ourApp = express()
ourApp.listen(3000)

ourApp.use(express.urlencoded({
  extended: false
}))


ourApp.get('/', function (req, res) {
    res.send(`
<form action="/answer" method="POST">
    <p>晴天的天空是什麼顏色？</p>
    <input name="skyColor" autocomplete="off">
    <button>送出答案</button>
  </form>
  `);
})


ourApp.post('/answer', function (req, res) {
	if (req.body.skyColor.toUpperCase() == "BLUE") {
		res.send(`
            <p>恭喜您，答對了。這是正確答案</p>
            <a href="/">回首頁</a>
        `)
    } else {
    	 res.send(`
           <p>真可惜，答錯了。</p>
           <a href="/">回首頁</a>
       `)
    }
})


ourApp.get('/answer', function (req, res) {
	res.send("迷路了嗎? 這裡什麼都沒有")
})

