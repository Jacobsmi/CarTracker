const express = require('express');
const { sequelize, User }= require('./database')



const app = express()
const port = 5000

sequelize
  .authenticate()
  .then(() => {
    console.log('Connection has been established successfully.');
  })
  .catch(err => {
    console.error('Unable to connect to the database:', err);
  });

app.get('/', (req, res) => {
  res.send('Hello World')
})

app.post('/createuser', async (req, res)=>{
  const newUser = User.build({ name: "Jane Doe", username: "janedoe", password: 'Pass1234!' });
  await newUser.save();
  res.send(JSON.stringify({"success":true}))
})

app.listen(port, () => {
  console.log("API running at http://localhost:5000")
})