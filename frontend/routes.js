const express = require('express');
const router = express.Router();
const axios = require('axios');
const multer = require('multer');
const { LocalStorage } = require('node-localstorage');
const localStorage = new LocalStorage('./scratch');

const storage = multer.diskStorage({
    destination: function (req, file, cb) {
        cb(null, 'uploads/')
    },
    filename: function (req, file, cb) {
        cb(null, file.originalname)
    }
});
const upload = multer({
    storage: storage,
    limits: { fileSize: 5 * 1024 * 1024 },
}).single('image');

router.get('/', (req, res) => {
    console.log("start / ")
    res.status(200).render('register');
});

router.post('/register', async (req, res)=>{
    console.log("start register")
    const { email, username, photo, aboutMe } = req.body;
    axios.post('http://localhost:8080/getvcode', { email })
        .then(response => {
            if (response.status == 200) {
                res.redirect(`/code?email=${email}&username=${username}&photo=${photo}&aboutMe=${aboutMe}`);
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});

router.get('/code', async(req, res)=>{
    console.log("start code get ");
    const { username, email, photo, aboutMe } = req.query;
    console.log(username, email, photo, aboutMe);
    res.render('code', { username, email, photo, aboutMe });
});

router.post('/code', async(req, res)=>{
    console.log("start code post");
    const { username, email } = req.query;
    console.log(username, email);
    const code = req.body.code;
    axios.post('http://localhost:8080/checkvcode', {
        email: email,
        code: code
    })
        .then(response => {
            if (response.status == 200) {
                res.redirect(`/getpass?email=${email}&username=${username}`);
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});

router.get('/getpass', (req, res)=>{
    console.log("start getpass get");
    const { username, email, photo, aboutMe } = req.query;
    res.render('getpass', { username, email, photo, aboutMe });
});

router.post('/getpass', (req, res)=>{
    console.log("start getpass post");
    const { username, email } = req.query;
    const password = req.body.password;
    axios.post('http://localhost:8080/register', {
        email: email,
        username: username,
        password: password,
        photo: req.body.photo,
        aboutMe: req.body.aboutMe
    })
        .then(response => {
            if (response.status == 200) {
                const token = response.data.token;
                localStorage.setItem('token', token);
                res.redirect(`/login`);
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});

router.get('/login', (req, res) => {
    console.log("start login get");
    res.status(200).render('login');
});

router.post('/login', async (req, res)=>{
    console.log("start login post");
    const { email, username, password } = req.body;
    axios.post('http://localhost:8080/login', { email, username, password })
        .then(response => {
            const token = response.data.token;
            localStorage.setItem('token', token);

            if (response.status == 200) {
                res.redirect(`index`);
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});
router.get('/index', (req, res) => {
    console.log("start index");

    const token = localStorage.getItem('token');
    let { filter, sort_by } = req.query;

    // Устанавливаем значения по умолчанию, если они не указаны в запросе
    filter = filter || ""; // Пустая строка, если filter не указан
    sort_by = sort_by || "name";

    axios.get(`http://localhost:8080/all?filter=${filter}&sort_by=${sort_by}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
        .then(response => {
            if (response.status == 200) {
                const sortBy = req.query.sortBy || 'defaultSortBy'; // Значение по умолчанию, если sortBy не указан
                res.render(`index`, {
                    books: response.data.books,
                    sortBy: sortBy
                });
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});
router.get('/book/:id/read', (req, res) => {
    const token = localStorage.getItem('token');
    const bookId = req.params.id;

    axios.get(`http://localhost:8080/book/${bookId}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
        .then(response => {
            if (response.status == 200) {
                const book = response.data.book;
                res.render('book_info', { book });
            } else {
                res.status(response.status).json({ error: "Failed to fetch book information" });
            }
        })
        .catch(error => {
            console.error('Error:', error);
            res.status(500).json({ error: "Internal Server Error" });
        });
});

router.post('/update', async (req, res) => {
    try {
        const token = req.headers.authorization;
        const userData = req.body; // Данные пользователя для обновления

        const response = await axios.post('http://localhost:8080/update', userData, {
            headers: {
                'Authorization': token
            }
        });

        if (response.status === 200) {
            const updatedUser = response.data.user;
            res.status(200).json({ message: "User updated successfully", user: updatedUser });
        } else {
            res.status(response.status).json({ error: "Failed to update user" });
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).json({ error: "Internal Server Error" });
    }
});

router.get('/profile', async (req, res) => {
    try {
        const token = req.headers.authorization;

        const response = await axios.get('http://localhost:8080/profile', {
            headers: {
                'Authorization': token
            }
        });

        if (response.status === 200) {
            const userProfile = response.data.user;
            res.status(200).json({ user: userProfile });
        } else {
            res.status(response.status).json({ error: "Failed to fetch user profile" });
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).json({ error: "Internal Server Error" });
    }
});
router.post('/books/:id/chapters', async (req, res) => {
    try {
        const token = req.headers.authorization;
        const bookId = req.params.id;
        const { name, words } = req.body;

        const response = await axios.post(`http://localhost:8080/books/${bookId}/chapters`, {
            name,
            words
        }, {
            headers: {
                'Authorization': token
            }
        });

        if (response.status === 200) {
            const newChapter = response.data.chapter;
            res.status(200).json({ chapter: newChapter });
        } else {
            res.status(response.status).json({ error: "Failed to add chapter" });
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).json({ error: "Internal Server Error" });
    }
});
router.post('/books/create', async (req, res) => {
    try {
        const token = req.headers.authorization;
        const { name, photo, briefInformation, genre, author, translator, finished } = req.body;

        const response = await axios.post('http://localhost:8080/books/create', {
            name,
            photo,
            briefInformation,
            genre,
            author,
            translator,
            finished
        }, {
            headers: {
                'Authorization': token
            }
        });

        if (response.status === 200) {
            const createdBook = response.data.book;
            res.status(200).json({ book: createdBook });
        } else {
            res.status(response.status).json({ error: "Failed to create book" });
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).json({ error: "Internal Server Error" });
    }
});
router.put('/books/:id/update', async (req, res) => {
    try {
        const token = req.headers.authorization;
        const bookId = req.params.id;
        const { name, photo, briefInformation, genre, author, translator, finished } = req.body;

        const response = await axios.put(`http://localhost:8080/books/${bookId}/update`, {
            name,
            photo,
            briefInformation,
            genre,
            author,
            translator,
            finished
        }, {
            headers: {
                'Authorization': token
            }
        });

        if (response.status === 200) {
            const updatedBook = response.data.book;
            res.status(200).json({ book: updatedBook });
        } else {
            res.status(response.status).json({ error: "Failed to update book" });
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).json({ error: "Internal Server Error" });
    }
});

router.delete('/books/:id/delete', async (req, res) => {
    try {
        const token = req.headers.authorization;
        const bookId = req.params.id;

        const response = await axios.delete(`http://localhost:8080/books/${bookId}/delete`, {
            headers: {
                'Authorization': token
            }
        });

        if (response.status === 200) {
            res.status(200).json({ message: "Book deleted successfully" });
        } else {
            res.status(response.status).json({ error: "Failed to delete book" });
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).json({ error: "Internal Server Error" });
    }
});
