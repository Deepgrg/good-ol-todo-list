import axios from "axios";

const instance = axios.create({
  baseURL: "https://good-ol-todo-list.herokuapp.com",
  timeout: 100000,
  headers: { "Content-Type": "application/json" },
});

export default instance;
