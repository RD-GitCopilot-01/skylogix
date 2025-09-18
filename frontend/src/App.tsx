import { useState, useEffect } from 'react'
import { Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'

interface Todo {
  id: number
  title: string
  completed: boolean
  created_at: string
}

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000'

function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [newTodo, setNewTodo] = useState('')
  const [loading, setLoading] = useState(false)

  const fetchTodos = async () => {
    try {
      const response = await fetch(`${API_URL}/api/todos`)
      if (response.ok) {
        const data = await response.json()
        setTodos(data || [])
      }
    } catch (error) {
      console.error('Error fetching todos:', error)
    }
  }

  const createTodo = async () => {
    if (!newTodo.trim()) return
    
    setLoading(true)
    try {
      const response = await fetch(`${API_URL}/api/todos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: newTodo.trim() }),
      })
      
      if (response.ok) {
        const todo = await response.json()
        setTodos(prev => [...prev, todo])
        setNewTodo('')
      }
    } catch (error) {
      console.error('Error creating todo:', error)
    } finally {
      setLoading(false)
    }
  }

  const updateTodo = async (id: number, completed: boolean) => {
    try {
      const response = await fetch(`${API_URL}/api/todos/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ completed }),
      })
      
      if (response.ok) {
        const updatedTodo = await response.json()
        setTodos(prev => prev.map(todo => 
          todo.id === id ? updatedTodo : todo
        ))
      }
    } catch (error) {
      console.error('Error updating todo:', error)
    }
  }

  const deleteTodo = async (id: number) => {
    try {
      const response = await fetch(`${API_URL}/api/todos/${id}`, {
        method: 'DELETE',
      })
      
      if (response.ok) {
        setTodos(prev => prev.filter(todo => todo.id !== id))
      }
    } catch (error) {
      console.error('Error deleting todo:', error)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      createTodo()
    }
  }

  useEffect(() => {
    fetchTodos()
  }, [])

  const completedCount = todos.filter(todo => todo.completed).length
  const totalCount = todos.length

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-2xl mx-auto px-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-3xl font-bold text-center text-gray-800">
              Todo App
            </CardTitle>
            <p className="text-center text-gray-600">
              {totalCount > 0 && (
                <span>
                  {completedCount} of {totalCount} completed
                </span>
              )}
            </p>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Input
                type="text"
                placeholder="Add a new todo..."
                value={newTodo}
                onChange={(e) => setNewTodo(e.target.value)}
                onKeyPress={handleKeyPress}
                className="flex-1"
              />
              <Button 
                onClick={createTodo} 
                disabled={loading || !newTodo.trim()}
                className="px-4"
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>

            <div className="space-y-2">
              {todos.length === 0 ? (
                <p className="text-center text-gray-500 py-8">
                  No todos yet. Add one above!
                </p>
              ) : (
                todos.map((todo) => (
                  <div
                    key={todo.id}
                    className={`flex items-center gap-3 p-3 rounded-lg border ${
                      todo.completed 
                        ? 'bg-gray-50 border-gray-200' 
                        : 'bg-white border-gray-300'
                    }`}
                  >
                    <Checkbox
                      checked={todo.completed}
                      onCheckedChange={(checked) => 
                        updateTodo(todo.id, checked as boolean)
                      }
                    />
                    <span
                      className={`flex-1 ${
                        todo.completed
                          ? 'line-through text-gray-500'
                          : 'text-gray-800'
                      }`}
                    >
                      {todo.title}
                    </span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => deleteTodo(todo.id)}
                      className="text-red-500 hover:text-red-700 hover:bg-red-50"
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default App
