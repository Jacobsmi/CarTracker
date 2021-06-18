import Image from 'next/image'
import styles from '../styles/Signup.module.css'

export default function Signup() {
    function submitForm() {
        console.log("Submitted")
    }
    return (
        <div className={styles.container}>
            <header className={styles.header}>
                Sign Up
            </header>
            <main className={styles.main}>
                <input type='text' id='name' className={styles.textInput} placeholder='Full Name' />

                <input type='text' id='username' className={styles.textInput} placeholder='Username' />

                <input type='password' id='password' className={styles.textInput} placeholder='Password' />

                <input type='password' id='confirm' className={styles.textInput} placeholder='Confirm Password' />

                <input type="button" value='Sign Up' className={styles.submit} id='submit' onClick={submitForm} />

            </main>
        </div>
    )
}