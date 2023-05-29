import { SiProbot } from "react-icons/si";
import { IoAccessibilityOutline } from "react-icons/io5";
import styles from "./_conversationBox.module.scss";
import { useEffect, useState } from "react";

interface ConvoBoxProps {
  from: string;
  inputValue: string;
}

export const ConversationBox = (props: ConvoBoxProps) => {
  const { from, inputValue } = props;
  const [message, setMessage] = useState(inputValue || "");

  useEffect(() => {
    const truncatedMessage =
      message.length > 250 ? message.slice(0, 250) + "..." : message;
    setMessage(truncatedMessage);
  }, [message]);

  return (
    <div className={styles.convoContainer}>
      {from == "john" ? <IoAccessibilityOutline size={"25px"} /> : null}
      <div className={styles.convoInput}>
        <div> {message}</div>
      </div>
      {from == "doe" ? <SiProbot size={"25px"} /> : null}
    </div>
  );
};
