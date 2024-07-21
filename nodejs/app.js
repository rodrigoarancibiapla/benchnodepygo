const express = require('express');
const app = express();

function generateEmail() {
    const characters = 'abcdefghijklmnopqrstuvwxyz';
    const usernameLength = 7;
    const domainLength = 5;

    const getRandomString = (length) => {
        return Array.from({ length }, () => characters.charAt(Math.floor(Math.random() * characters.length))).join('');
    };

    const username = getRandomString(usernameLength);
    const domain = getRandomString(domainLength);
    
    return `${username}@${domain}.com`;
}

function generatePerson() {
    const email = generateEmail();
    const age = Math.floor(Math.random() * 100) + 1;
    return { email, age };
}
app.set('json spaces',2)
app.get('/people', (req, res) => {
    const people = Array.from({ length: 5000 }, generatePerson);
    res.header('Content-Type','application/json')
    res.json(people.slice(0, 20));
});

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});
