const express = require('express');
const { sequelize, User }= require('./database')
const bodyParser = require('body-parser')



const app = express()
app.use(express.json())
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
  const newUser = User.build({ name: req.body.name, username: req.body.username, password: req.body.password });
  try{
    await newUser.save()
  }catch (err) {
    if (err.name === "SequelizeUniqueConstraintError"){
      res.send(JSON.stringify({"error": "account_already_exists"}))
    }
  }

  res.send(JSON.stringify({"success":true}))
})

app.listen(port, () => {
  console.log("API running at http://localhost:5000")
})