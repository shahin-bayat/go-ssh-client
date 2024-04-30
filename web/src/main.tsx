import React from 'react'
import ReactDOM from 'react-dom/client'
import {createBrowserRouter, RouterProvider} from "react-router-dom"

import Login from "./routes/Login.tsx"
import Admin from "./routes/Admin.tsx";

const router = createBrowserRouter([
    {path: "/", element: <Login/>},
    {path: "/admin", element: <Admin/>}
])


ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <RouterProvider router={router}/>
    </React.StrictMode>,
)
