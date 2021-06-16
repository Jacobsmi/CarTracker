import Image from 'next/image'
import styles from '../styles/Signup.module.css'

export default function Signup(){
    function submitForm(){
        console.log("Submitted")
    }
    return(
        <div className={styles.container}>
            <main className={styles.form}>
                <h1>Sign Up</h1>
                <form>
                    <input type='text' id='username' className={styles.textInput} placeholder='Username' /> <br />
                    
                    <input type='password' id='password' className={styles.textInput} placeholder='Password' /> <br />
                    
                    <input type='password' id='confirm' className={styles.textInput} placeholder='Confirm Password' /> <br />

                    <input type="button" value='Sign Up' className={styles.submit} id='submit' style={{marginTop: '5vh'}} onClick={submitForm} />
                </form>
            </main>
            <main className={styles.image}>
                <Image src='/signup.svg' layout='fill' />
            </main>
        </div>
    )
}