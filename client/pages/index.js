import axios from "../utils/axios";
import { useState, useEffect } from "react";

export default function Home() {
  const [todo, setTodo] = useState("");
  const [todos, setTodos] = useState([]);
  const [loading, setLoading] = useState(false);
  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    setLoading(true);
    try {
      const result = await axios.get("/");
      setLoading(false);
      setTodos(result.data.Data);
    } catch (err) {
      setLoading(false);
      console.log(err);
    }
  };

  const createTodo = async () => {
    setLoading(true);
    try {
      await axios.post("/", {
        title: todo,
        completed: false,
      });
      setTodo("");
      fetchTodos();
    } catch (err) {
      setLoading(false);
      console.log(err);
    }
  };

  const updateTodo = async (todoId, todoStatus) => {
    setLoading(true);

    try {
      await axios.put(`/${todoId}`);
      fetchTodos();
    } catch (error) {
      setLoading(false);
      console.log(error);
    }
  };

  const deleteTodo = async (todoId) => {
    setLoading(true);
    try {
      await axios.delete(`/${todoId}`);
      fetchTodos();
    } catch (error) {
      setLoading(false);
      console.log(error);
    }
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className="w-full  h-14 flex items-center justify-center">
        <p className="font-semibold text-3xl underline underline-offset-2">
          Good ol&apos; Todo List
        </p>
      </div>
      <div className="w-full flex items-center justify-center">
        <div className="w-3/6 p-4 border border-gray-400 rounded-md shadow-md">
          <div className="w-full text-center font-semibold text-xl underline underline-offset-1 mb-5">
            Add new todo
          </div>
          <div className="w-full ">
            <form
              className="w-full flex items-center justify-between  px-3 pb-5"
              onSubmit={createTodo}
            >
              <label>Todo title</label>
              <input
                className="ml-5 outline-none border border-gray-300 px-2 py-2 rounded-md"
                type="text"
                value={todo}
                onChange={(e) => setTodo(e.target.value)}
                placeholder="Add a title to your todo"
              />
              <button
                className="px-2 py-2 bg-green-400 rounded-md"
                type="submit"
              >
                Create
              </button>
            </form>
          </div>
          <div className="w-full text-center font-semibold text-xl underline underline-offset-1 mb-5">
            Todos
          </div>

          {!loading &&
            todos.map((element) => {
              return (
                <div
                  key={element.todo_id}
                  className="w-full flex items-center justify-between px-3 py-3"
                >
                  <p>{element.title}</p>
                  <div className="space-x-3">
                    <input
                      type="checkbox"
                      defaultChecked={element.completed}
                      onChange={() =>
                        updateTodo(element.todo_id, !element.completed)
                      }
                    />
                    <button
                      className="px-2 py-2 bg-red-400 rounded-md"
                      onClick={() => deleteTodo(element.todo_id)}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              );
            })}
        </div>
      </div>
    </div>
  );
}
