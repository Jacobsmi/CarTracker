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
      <header className={styles.header}>
        <div className={styles.logo}>
          Car Tracker
        </div>
        <div className={styles.links}>
          <Link href='/login'>
            <a>Login</a>
          </Link>
          <Link href='/signup'>
            <a id={styles.signupLink}>Sign Up</a>
          </Link>
        </div>
      </header>

      <main className={styles.body}>
        <div className={styles.head1}>
          Never lose the cars of your dreams
        </div>
        <div className={styles.head2}>
          Track cars you want and cars you have owned so you never forget
        </div>
        <Link href='/signup'>
          <a className={styles.getstartedButton}>Get Started</a>
        </Link>
      </main>
    </div>
  )
}
