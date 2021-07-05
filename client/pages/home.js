import { useEffect } from 'react';
import Router from 'next/router'

export default function Home(){
  useEffect(()=>{
    async function checkToken() {
      const resp = await fetch("http://localhost:5000/checktoken",{
        method: "POST",
        credentials: "include",
        headers: {
          'Content-Type': 'application/json'
        }
      })
      const respJSON = await resp.json()
      if (respJSON.Success != true) {
        Router.push('/')
      }
    }
    checkToken()
  }, [])
  return(
    <div>
      Home Page
    </div>
  )
}