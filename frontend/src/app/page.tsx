import styles from './page.module.css'
import MyTable from "./ui/table";


export default function Home() {
  return (
    <main className={styles.main}>
      <MyTable></MyTable>
    </main>
  )
}
