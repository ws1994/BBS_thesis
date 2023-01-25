let express = require('express');
let app = express();

app.get('/', function (req, res) {
    res.render('index', {
        'title': '首頁',
        'titleH2': '<h2>第二級標題</h2>',
        'show': true,
        'foods': ['apple', 'banana', 'mongo']
    });
})


let port = 3000;
app.listen(port);

let engine = require('ejs-locals');
app.engine('ejs', engine);
app.set('views', '../views');
app.set('view engine', 'ejs');

app.get('/user', function (req, res) {
    res.render('user', {
        'title': '用户'
    });
})