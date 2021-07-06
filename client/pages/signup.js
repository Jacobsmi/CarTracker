import { useState } from 'react'
import styles from '../styles/Signup.module.css'
import Router from 'next/router'

export default function Signup() {
  // Create a state to track if there are errors in the form or not
  const [error, setError] = useState(false)

  async function submitForm() {
    // Create a string that is generated that tracks errors in the form
    let errorString = 'Error(s):<ul>'

    // Get and validate name
    const name = document.getElementById("name").value
    const validName = /^[a-zA-Z-']+ [a-zA-Z-']+$/.test(name)

    // If name not valid add to error string
    if (!validName) {
      errorString += '<li>Invalid Name</li>'
    }

    // Get and validate username
    const username = document.getElementById('username').value
    const validUsername = /^[a-zA-Z0-9]+$/.test(username)

    // If not valid add to error string
    if (!validUsername) {
      setError(true)
      errorString += '<li>Invalid Username</li>'
    }

    // Get and validate password
    const password = document.getElementById('password').value
    const validPassword = /^(?=.*[0-9])(?=.*[!@#$%^&*])(?=.*[A-Z])[a-zA-Z0-9!@#$%^&*]{8,}$/.test(password)

    // If not valid add to error string
    if (!validPassword) {
      errorString += '<li>Invalid Password - Must be 8 digits, have a capital letter, one digit, and one special character</li>'
    }

    //Ensure passwords are matching
    const confirmPass = document.getElementById('confirm').value
    const passMatch = (confirmPass == password)

    //If not add to error string
    if (!passMatch) {
      errorString += '<li>Passwords do not match</li>'
    }

    // If all fields are valid do an API call
    if (validName && validUsername && validPassword && passMatch) {
      setError(false)
      const result = await fetch('http://localhost:5000/signup', {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
          "name": name,
          "username": username,
          "password": password
        })
      })
      const resultJSON = await result.json()
      console.log(resultJSON)
      if (resultJSON.Success != true) {
        setError(true)
        if (resultJSON.Err == "duplicate_user") {
          errorString += '<li>Username already exists</li>'
        } else if (resultJSON.Err == "unhandled_db_error") {
          errorString += '<li>Unhandled DB Error <br> Please Check Database or Contact Admin</li>'
        } else {
          errorString += '<li>Unhandled Error</li>'
        }
      }else if(resultJSON.Success == true){
        Router.push('/home')
      }
    }
    // Otherwise display errors on the page
    else {
      setError(true)
    }
    // Close the unordered list in the error string
    errorString += '</ul>'
    document.getElementById('errors').innerHTML = errorString
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        Sign Up
            </header>
      <main className={styles.main}>
        <div className={styles.errors} id='errors' style={{ 'display': `${error ? 'block' : 'none'}` }}></div>
        <input type='text' id='name' className={styles.textInput} placeholder='Full Name' />

        <input type='text' id='username' className={styles.textInput} placeholder='Username' />

        <input type='password' id='password' className={styles.textInput} placeholder='Password' />

        <input type='password' id='confirm' className={styles.textInput} placeholder='Confirm Password' />

        <input type="button" value='Sign Up' className={styles.submit} id='submit' onClick={submitForm} />

      </main>
    </div>
  )
}