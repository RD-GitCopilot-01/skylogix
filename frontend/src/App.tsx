import { useState, useEffect } from 'react'
import { Plus, Trash2, CheckCircle2, Circle, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Progress } from '@/components/ui/progress'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'

interface Todo {
  id: number
  title: string
  completed: boolean
  created_at: string
}

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000'

function App() {
  const [todos, setTodos] = useState<Todo[]>([
    { id: 1, title: "Complete project documentation", completed: false, created_at: "2024-01-15T10:00:00Z" },
    { id: 2, title: "Review code changes", completed: true, created_at: "2024-01-15T09:30:00Z" },
    { id: 3, title: "Update dependencies", completed: false, created_at: "2024-01-15T09:00:00Z" },
    { id: 4, title: "Fix responsive design issues", completed: true, created_at: "2024-01-15T08:30:00Z" },
    { id: 5, title: "Implement user authentication", completed: false, created_at: "2024-01-15T08:00:00Z" },
    { id: 6, title: "Write unit tests", completed: false, created_at: "2024-01-15T07:30:00Z" },
    { id: 7, title: "Optimize database queries", completed: true, created_at: "2024-01-15T07:00:00Z" },
    { id: 8, title: "Deploy to staging environment", completed: false, created_at: "2024-01-15T06:30:00Z" },
    { id: 9, title: "Create API documentation", completed: false, created_at: "2024-01-15T06:00:00Z" },
    { id: 10, title: "Set up monitoring alerts", completed: true, created_at: "2024-01-15T05:30:00Z" },
    { id: 11, title: "Refactor legacy components", completed: false, created_at: "2024-01-15T05:00:00Z" },
    { id: 12, title: "Update user interface design", completed: false, created_at: "2024-01-15T04:30:00Z" },
    { id: 13, title: "Implement caching strategy", completed: true, created_at: "2024-01-15T04:00:00Z" },
    { id: 14, title: "Configure CI/CD pipeline", completed: false, created_at: "2024-01-15T03:30:00Z" },
    { id: 15, title: "Conduct security audit", completed: false, created_at: "2024-01-15T03:00:00Z" },
    { id: 16, title: "Optimize image loading", completed: true, created_at: "2024-01-15T02:30:00Z" },
    { id: 17, title: "Add error handling", completed: false, created_at: "2024-01-15T02:00:00Z" },
    { id: 18, title: "Update README file", completed: false, created_at: "2024-01-15T01:30:00Z" },
    { id: 19, title: "Implement dark mode", completed: true, created_at: "2024-01-15T01:00:00Z" },
    { id: 20, title: "Performance testing", completed: false, created_at: "2024-01-15T00:30:00Z" },
    { id: 21, title: "Implement search functionality", completed: false, created_at: "2024-01-15T00:00:00Z" },
    { id: 22, title: "Add data validation", completed: true, created_at: "2024-01-14T23:30:00Z" },
    { id: 23, title: "Create user onboarding flow", completed: false, created_at: "2024-01-14T23:00:00Z" },
    { id: 24, title: "Optimize bundle size", completed: false, created_at: "2024-01-14T22:30:00Z" },
    { id: 25, title: "Add keyboard shortcuts", completed: true, created_at: "2024-01-14T22:00:00Z" },
    { id: 26, title: "Implement drag and drop", completed: false, created_at: "2024-01-14T21:30:00Z" },
    { id: 27, title: "Add export functionality", completed: false, created_at: "2024-01-14T21:00:00Z" },
    { id: 28, title: "Create admin dashboard", completed: true, created_at: "2024-01-14T20:30:00Z" },
    { id: 29, title: "Implement real-time notifications", completed: false, created_at: "2024-01-14T20:00:00Z" },
    { id: 30, title: "Add multi-language support", completed: false, created_at: "2024-01-14T19:30:00Z" }
  ])
  const [newTodo, setNewTodo] = useState('')
  const [loading, setLoading] = useState(false)
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [completingId, setCompletingId] = useState<number | null>(null)

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
    setCompletingId(id)
    
    setTodos(prev => prev.map(todo => 
      todo.id === id ? { ...todo, completed } : todo
    ))
    
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
      } else {
        setTodos(prev => prev.map(todo => 
          todo.id === id ? { ...todo, completed: !completed } : todo
        ))
      }
    } catch (error) {
      console.error('Error updating todo:', error)
      setTodos(prev => prev.map(todo => 
        todo.id === id ? { ...todo, completed: !completed } : todo
      ))
    } finally {
      setCompletingId(null)
    }
  }

  const deleteTodo = async (id: number) => {
    setDeletingId(id)
    
    try {
      const response = await fetch(`${API_URL}/api/todos/${id}`, {
        method: 'DELETE',
      })
      
      if (response.ok) {
        setTodos(prev => prev.filter(todo => todo.id !== id))
      }
    } catch (error) {
      console.error('Error deleting todo:', error)
    } finally {
      setDeletingId(null)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      createTodo()
    } else if (e.key === 'Escape') {
      setNewTodo('')
    }
  }

  useEffect(() => {
    fetchTodos()
  }, [])

  const completedCount = todos.filter(todo => todo.completed).length
  const totalCount = todos.length
  const progressPercentage = totalCount > 0 ? (completedCount / totalCount) * 100 : 0

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 py-4 px-4">
      <div className="w-full">
        <Card className="shadow-xl border-0 bg-white/80 backdrop-blur-sm">
          <CardHeader className="text-center space-y-4 pb-6">
            <div className="flex items-center justify-center gap-2 mb-2">
              <CheckCircle2 className="w-8 h-8 text-indigo-600" />
              <CardTitle className="text-4xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                Todo App
              </CardTitle>
            </div>
            
            {totalCount > 0 && (
              <div className="space-y-3">
                <div className="flex items-center justify-center gap-2 text-sm text-gray-600">
                  <Circle className="w-4 h-4" />
                  <span className="font-medium">
                    {completedCount} of {totalCount} completed ({Math.round(progressPercentage)}%)
                  </span>
                </div>
                <div className="max-w-xs mx-auto">
                  <Progress 
                    value={progressPercentage} 
                    className="h-2 bg-gray-200"
                  />
                </div>
              </div>
            )}
          </CardHeader>
          
          <CardContent className="space-y-6">
            <div className="flex gap-2">
              <Input
                type="text"
                placeholder="Add a new todo... (Press Enter to add, Esc to clear)"
                value={newTodo}
                onChange={(e) => setNewTodo(e.target.value)}
                onKeyPress={handleKeyPress}
                className="flex-1 border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 transition-colors"
                disabled={loading}
              />
              <Button 
                onClick={createTodo} 
                disabled={loading || !newTodo.trim()}
                className="px-4 bg-indigo-600 hover:bg-indigo-700 transition-all duration-200 shadow-md hover:shadow-lg"
              >
                {loading ? (
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                ) : (
                  <Plus className="w-4 h-4" />
                )}
              </Button>
            </div>

            <div className="space-y-3">
              {todos.length === 0 ? (
                <div className="text-center py-12 space-y-4">
                  <div className="w-16 h-16 mx-auto bg-gradient-to-br from-indigo-100 to-purple-100 rounded-full flex items-center justify-center">
                    <AlertCircle className="w-8 h-8 text-indigo-500" />
                  </div>
                  <div>
                    <p className="text-gray-600 font-medium">No todos yet!</p>
                    <p className="text-gray-500 text-sm">Add your first todo above to get started</p>
                  </div>
                </div>
              ) : (
                <div className="space-y-2">
                  {todos.map((todo, index) => {
                    const colorSchemes = [
                      'bg-gradient-to-r from-red-50 to-pink-50 border-red-200 hover:border-red-300',
                      'bg-gradient-to-r from-orange-50 to-yellow-50 border-orange-200 hover:border-orange-300',
                      'bg-gradient-to-r from-green-50 to-emerald-50 border-green-200 hover:border-green-300',
                      'bg-gradient-to-r from-blue-50 to-cyan-50 border-blue-200 hover:border-blue-300',
                      'bg-gradient-to-r from-purple-50 to-violet-50 border-purple-200 hover:border-purple-300',
                      'bg-gradient-to-r from-pink-50 to-rose-50 border-pink-200 hover:border-pink-300'
                    ];
                    const colorScheme = todo.completed 
                      ? 'bg-gradient-to-r from-green-50 to-emerald-50 border-green-200 hover:border-green-300'
                      : colorSchemes[index % colorSchemes.length];
                    
                    return (
                    <div
                      key={todo.id}
                      className={`group flex items-center gap-3 p-4 rounded-xl border transition-all duration-300 hover:shadow-md ${colorScheme} ${completingId === todo.id ? 'animate-pulse' : ''} ${
                        deletingId === todo.id ? 'opacity-50 scale-95' : 'animate-in slide-in-from-top-2'
                      }`}
                      style={{ animationDelay: `${index * 50}ms` }}
                    >
                      <Checkbox
                        checked={todo.completed}
                        onCheckedChange={(checked) => 
                          updateTodo(todo.id, checked as boolean)
                        }
                        disabled={completingId === todo.id}
                        className="data-[state=checked]:bg-green-600 data-[state=checked]:border-green-600"
                      />
                      <span
                        className={`flex-1 transition-all duration-300 ${
                          todo.completed
                            ? 'line-through text-gray-500 font-medium'
                            : 'text-gray-800 group-hover:text-indigo-700'
                        }`}
                      >
                        {todo.title}
                      </span>
                      
                      <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button
                            variant="ghost"
                            size="sm"
                            disabled={deletingId === todo.id}
                            className="opacity-0 group-hover:opacity-100 transition-all duration-200 text-red-500 hover:text-red-700 hover:bg-red-50 hover:scale-105"
                          >
                            {deletingId === todo.id ? (
                              <div className="w-4 h-4 border-2 border-red-500 border-t-transparent rounded-full animate-spin" />
                            ) : (
                              <Trash2 className="w-4 h-4" />
                            )}
                          </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Delete Todo</AlertDialogTitle>
                            <AlertDialogDescription>
                              Are you sure you want to delete "{todo.title}"? This action cannot be undone.
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction
                              onClick={() => deleteTodo(todo.id)}
                              className="bg-red-600 hover:bg-red-700"
                            >
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </div>
                    );
                  })}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default App
