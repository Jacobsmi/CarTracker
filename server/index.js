const express = require('express');
const { Sequelize, Model, DataTypes } = require('sequelize');
require('dotenv').config()


const app = express()
const sequelize = new Sequelize(`postgres://${process.env.DB_USER}:${process.env.DB_PASS}@${process.env.DB_HOST}:${process.env.DB_PORT}/${process.env.DB_NAME}`);
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

app.listen(port, () => {
  console.log("API running at http://localhost:5000")
})