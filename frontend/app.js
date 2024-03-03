const express = require('express');
const bodyParser = require('body-parser');
const session = require('express-session');
const path = require('path');
const routes = require('./routes/routes');
const cookieParser = require('cookie-parser');

const app = express();
const port = 3000;
app.set('view engine', 'ejs');
app.use('/uploads', express.static('uploads'));
app.use(express.static('public'));

app.use(cookieParser());
app.use(session({
    secret: 'secret',
    resave: true,
    saveUninitialized: true
}));
app.use(express.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(express.static(path.join(__dirname, 'public')));


app.use('/', routes);
app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});