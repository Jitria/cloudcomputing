import {useState} from "react";

export default function Hello(props){
    const [name, setName] = useState('Mike');
    const [age, setAge] = useState(props.age)

    return (
    <div>
        <h1>state</h1>
        <h2 id="name">{name}({age})</h2>
        <button onClick={() => {
            setName(name === "Mike" ? "Jane" : "Mike");
            setAge(age+1);
            }}
        >
        Change
        </button>
    </div>
    );
}