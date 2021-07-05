import { useEffect } from 'react';

export default function Home(){
  useEffect(()=>{
    console.log(document.cookie)
  })
  return(
    <div>
      Home Page
    </div>
  )
}