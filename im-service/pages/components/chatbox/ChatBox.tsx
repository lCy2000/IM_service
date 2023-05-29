import styles from "./_chatbox.module.scss";
import { BsFillCaretRightFill } from "react-icons/bs";
import { ConversationBox } from "../conversationbox/ConversationBox";
import { useEffect, useState } from "react";

export const ChatBox = () => {
  const [response, setResponse] = useState(null);
  const [message, setMessage] = useState<string>("");
  const [allMessage, setAllMessage] = useState<any>();
  const [user, setUser] = useState("john");

  const sendMsgReq = async () => {
    try {
      const requestBody = {
        chat: "john:doe",
        text: message,
        sender: user,
      };
      const response = await fetch("/api/send", {
        method: "POST",
        body: JSON.stringify(requestBody),
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      setResponse(responseData);
      setUser(user == "john" ? "doe" : "john");
      setMessage("");
    } catch (error) {
      console.error("Error making POST request:", error);
    }
  };

  //NOTE: Instead of fetching 10 everytime, we should fetch the next one everytime which keep the prev msg there
  const pullMsgReq = async () => {
    try {
      const requestParams = {
        chat: "john:doe",
        cursor: "0",
        limit: "10",
        reverse: "false",
      };

      const queryParams = new URLSearchParams(requestParams).toString();

      const response = await fetch(
        `http://localhost:8080/api/pull/messages?${queryParams}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      const data = await response.json();
      setAllMessage(data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  useEffect(() => {
    pullMsgReq();
  }, [user]);

  return (
    <div className={styles.wrapperBox}>
      <div className={styles.chatBox}>
        <div className={styles.containerBox}>
          {allMessage &&
            allMessage?.messages?.map((message: any) => {
              return (
                <>
                  <ConversationBox
                    key={message.sender}
                    from={message.sender}
                    inputValue={message.text}
                  />
                </>
              );
            })}
        </div>
        <div className={styles.inputBox}>
          <input
            type="text"
            placeholder=""
            onChange={(e: any) => setMessage(e.target.value)}
            value={message}
          />
          <div
            className={styles.iconBox}
            onClick={() => {
              sendMsgReq();
            }}
          >
            <BsFillCaretRightFill size={"30px"} />
          </div>
        </div>
      </div>
    </div>
  );
};
