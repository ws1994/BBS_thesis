let express = require('express');
let engine = require('ejs-locals');
let app = express();
const bodyParser = require("body-parser")
let searchFile = require("../../magnetocorp/application/searchAllFileClass");
let addFile = require("../../digibank/application/initFileClass");
let listOffState = require("../../magnetocorp/application/listOffStateClassMag");
let requestFile = require("../../magnetocorp/application/requestFileClass");
let chainOfCustody = require("../../digibank/application/chainOfCustodyClass");

app.use(bodyParser.urlencoded({
  extended:true
}));

app.engine('ejs', engine);
app.set('views', '../views');
app.set('view engine', 'ejs');
app.use(express.json());

let port = 3000;
app.listen(port);
app.use(express.static('../public'

var requestAllFile = false;
var fileResArr = [];
app.get('/', function (req, res) {
if (!requestAllFile){
    searchFile.searchAllFile().then(function(fileRes){
        fileResArr = fileRes;
    });
    // requestAllFile = true;
}
  res.render('index', {
    'title': '链上文件信息列表',
    'fileRes':fileResArr
  });
})

app.post('/requestFile', function (req, res) { 
    requestFile.requestFile(req.body.filename);  
    console.log(req.body.filename);
    res.redirect('/');
})


// add data
app.get('/product', function (req, res) {
res.render('product', {
    'title': '数据发布',
    'subtitle': '请填写文件信息',
 });
})

app.post('/addFile', function (req, res) {    
    console.log(req.body);

    addFile.addFile(req.body.fileName, req.body.fileDes, req.body.fileRule);

    res.redirect('/contact'); 
})


var txChainArr = [];
app.get('/about', function (req, res) {
    res.render('about', {
      'title': '查看证据链',
      'txChainArr':txChainArr
 });
})

app.post('/chainOfCustody', function (req, res) {    
    console.log(req.body.filename);
    chainOfCustody.buildChainCustody(req.body.filename).then(function(txChain){
        console.log(txChain);
        txChainArr = txChain;
    });
    res.redirect("/about")
})


// off-state data list
var flag = false;
var dataArr = [];
app.get('/contact', function (req, res) {
    if(!flag){

        listOffState.listOffState().then(function(offstateList){
            console.log(offstateList)
            var ruleArr = offstateList.split("\n");
            console.log(ruleArr.length)
            
            for (i=0;i<(ruleArr.length-1);i++){
                console.log(ruleArr[i]);
                dataArr[i]=ruleArr[i];
            }
            console.log(dataArr);
        })
        // flag = true;
    }

    res.render('contact', {
        'title': '我的off-state数据',
        'dataArr':dataArr
    });
})

app.post('/download', function (req, res) {    
    console.log(req.body.filename);
    res.redirect('/contact'); 
})