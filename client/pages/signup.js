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
                    <label htmlFor='name'>Name</label>
                    <input type='text' id='name' /> <br />
                    <label htmlFor='email'>E-Mail</label>
                    <input type='text' id='email' /> <br />
                    <label htmlFor='password'>Password</label>
                    <input type='password' id='password' /> <br />
                    <label htmlFor='confirm'>Confirm Password</label>
                    <input type='password' id='confirm' /> <br />
                    <input type="button" value='Sign Up' id='submit' style={{marginTop: '5vh'}} onClick={submitForm} />
                </form>
            </main>
            <main className={styles.image}>
                <Image src='/signup.svg' layout='fill' />
            </main>
        </div>
    )
}