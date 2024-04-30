import {useState} from 'react'
import {Link} from "react-router-dom";

function Login() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')


    const onSubmit: React.FormEventHandler = async (e) => {
        e.preventDefault()
        const formData = new FormData()
        formData.append('username', username)
        formData.append('password', password)

        const response = await fetch("/login", {
            method: "POST",
            body: formData
        })
        const data = await response.json()
        console.log(data)
    }

    return (
        <>
            <form onSubmit={onSubmit}>
                <input type="text" value={username} onChange={(e) => setUsername(e.target.value)}/>
                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}/>
                <button type="submit">Submit</button>
            </form>
            <Link to="/admin">Admin</Link>
        </>
    )
}

export default Login
