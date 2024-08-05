const express = require('express');
const app = express();
const cluster = require('cluster');
const os = require('os');


const PORT = process.env.PORT || 3001;


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

const numWorkers = 4;//process.argv[2] || os.cpus().length;

if (cluster.isMaster) {
    console.log(`Master ${process.pid} is running`);

    // Fork workers.
    for (let i = 0; i < numWorkers; i++) {
        cluster.fork();
    }

    cluster.on('exit', (worker, code, signal) => {
        console.log(`Worker ${worker.process.pid} died`);
        // Optionally, you can restart the worker here if necessary
        // cluster.fork();
    });
} else {
    // Workers can share any TCP connection.
    // In this case, it is an HTTP server.
 
    app.listen(PORT, () => {
        console.log(`Worker ${process.pid} is running on port ${PORT}`);
    });

    // Your existing app code here
    // app.get('/', (req, res) => res.send('Hello World!'));
}