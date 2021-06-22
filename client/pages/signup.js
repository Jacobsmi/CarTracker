import Image from 'next/image'
import { useState } from 'react'
import styles from '../styles/Signup.module.css'

export default function Signup() {
  
  const [error, setError] = useState(false)
  
  async function submitForm() {
    let errorString = 'Error(s):<ul>'

    const name = document.getElementById("name").value
    const validName = /^[a-zA-Z-']+ [a-zA-Z-']+$/.test(name)
    
    if (!validName) {
      setError(true)
      errorString += '<li>Invalid Name</li>'
    }

    const username = document.getElementById('username').value
    const validUsername = /^[a-zA-Z0-9]+$/.test(username)

    if(!validUsername){
      errorString += '<li>Invalid Username</li>'
    }

    const password = document.getElementById('password').value
    const validPassword = /^(?=.*[0-9])(?=.*[!@#$%^&*])(?=.*[A-Z])[a-zA-Z0-9!@#$%^&*]{8,}$/.test(password)

    if(!validPassword){
      errorString += '<li>Invalid Password - Must be 8 digits, have a capital letter, one digit, and one special character</li>'
    }
    
    const confirmPass = document.getElementById('confirm').value
    const passMatch = (confirmPass == password)

    if(!passMatch){
      errorString += '<li>Passwords do not match</li>'
    }
    errorString += '</ul>'

    if(validName && validUsername && validPassword && passMatch){
      setError(false)
      const result = await fetch('/api/signup')
      const resultJSON = await result.json()
      console.log(resultJSON)
    }else{
      setError(true)
      document.getElementById('errors').innerHTML = errorString
    }
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