import type { NextPage } from "next";
import Head from "next/head";
import styles from "../styles/Home.module.css";
import { ChatBox } from "./components/chatbox/chatbox";

const Home: NextPage = () => {
  return (
    <div className={styles.container}>
      <Head>
        <title>Evercard</title>
        <meta name="description" content="Evercard" />
      </Head>

      <main>
        <ChatBox />
      </main>
    </div>
  );
};

export default Home;
