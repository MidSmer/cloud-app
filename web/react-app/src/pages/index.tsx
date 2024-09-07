import Head from 'next/head'
import Game from '@/components/game';

export default function Home() {
  return (
    <>
      <Head>
        <title>Home</title>
      </Head>
      <main className='min-h-screen bg-black text-white w-full'>
      <Game />
      </main>      
    </>
  );
}
