import axios from 'axios';
import { useCallback, useEffect, useState } from 'react'
import { HTTPURL } from "@/env";

interface Task {
  id: string;
  title: string;
  status: number;
  createdAt: string;
  updatedAt: string;
}

interface TasksResponse {
  data: {
    tasks: Task[]
  }
}

function App() {
  const [title, setTitle] = useState("");
  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    getTasks();
  }, []);

  const getTasks = useCallback(() => {
    axios.get(`${HTTPURL}api/task`).then((r: TasksResponse) => {
      setTasks(r.data.tasks);
    });
  }, []);

  const onAddClick = useCallback(() => {
    if (!title.match("^[a-zA-Z0-9]{1,10}$")) {
      return;
    }

    axios.post(`${HTTPURL}api/task`, { title }).then(() => {
      setTitle("");
      getTasks();
    });
  }, [title]);

  const onReloadClick = useCallback(() => {
    getTasks();
  }, []);

  const changeStatus = useCallback((id: string, status: string) => {
    axios.put(`${HTTPURL}api/task/${id}`, { status: Number(status) }).then(() => getTasks());
  }, []);

  const onRemoveClick = useCallback((id: string) => {
    axios.delete(`${HTTPURL}api/task/${id}`).then(() => getTasks());
  }, []);

  return (
    <>
      <div className='flex flex-col bg-stone-800 min-w-md h-screen p-2'>
        <div className='flex'>
          <input className='m-2 p-2 flex-1 rounded-lg bg-stone-700' type="text" placeholder='title...' maxLength={10} onChange={e => setTitle(e.target.value)} value={title} />
          <button className='m-2 p-2 hover:bg-stone-900 active:bg-stone-700 rounded-lg' onClick={onAddClick}>add</button>
          <button className='m-2 p-2 hover:bg-stone-900 active:bg-stone-700 rounded-lg' onClick={onReloadClick}>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className="size-5">
              <path d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
          </button>
        </div>
        <div className='flex-1 overflow-y-auto flex flex-col'>
          {
            tasks.map(task =>
              <div key={task.id} className='m-2 bg-stone-900 rounded-xl h-24 my-6 p-2 relative flex flex-col '>
                <p className='text-xs absolute -top-5'>ID: {task.id}</p>
                <button className='p-1 absolute right-0 -top-6 hover:bg-stone-900 active:bg-stone-700 rounded-sm' onClick={() => onRemoveClick(task.id)}>
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className="size-4">
                    <path d="M6 18 18 6M6 6l12 12" />
                  </svg>
                </button>
                <div className='flex flex-1 items-center'>
                  <span className='flex-1'>{task.title}</span>
                  <select className='bg-stone-900 flex-1' value={task.status} onChange={e => changeStatus(task.id, e.target.value)}>
                    <option value="0">Pending</option>
                    <option value="1">In progress</option>
                    <option value="2">Complete</option>
                  </select>
                </div>
                <div className='flex flex-1 items-center text-xs text-gray-400'>
                  <div className='flex-1'>
                    <p>Created At:</p>
                    <p>{task.createdAt}</p>
                  </div>
                  <div className='flex-1'>
                    <p>Updated At:</p>
                    <p>{task.updatedAt}</p>
                  </div>
                </div>
              </div>)
          }
        </div>
      </div>
    </>
  )
}

export default App
