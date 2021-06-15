import Head from 'next/head'
import Link from 'next/link'
import Image from 'next/image'
import styles from '../styles/Home.module.css'

export default function Home() {
  return (
    <div className={styles.container}>
      <Head>
        <title>Car Tracker-Welcome</title>
        <meta name="description" content="Landing page for Car Tracker" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <header className={styles.logo}>
        Car Tracker
      </header>
      <header className={styles.links}>
        <Link href='/signup'>
          <a>Login</a>
        </Link>
        <Link href='/signup'>
          <a id={styles.signupLink}>Sign Up</a>
        </Link>
      </header>

      <main className={styles.body}>
        <div className={styles.bodyText}>
          <h1>Track all the cool and exciting cars you have ever owned and want to own in the future</h1>
          <h2>Sign up and never forget a cool car ever again</h2>
          <Link href='/signup'>
            <a className={styles.getStartedLink}>Get Started</a>
          </Link>
        </div>
        <div className={styles.bodyImage}>
          <Image src='/car.svg' layout='fill' />
        </div>
      </main>
    </div>
  )
}
