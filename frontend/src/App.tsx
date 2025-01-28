import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ThemeProvider, createTheme } from '@mui/material'
import Layout from './components/layout/Layout'
import Dashboard from './pages/Dashboard'
import Login from './pages/Login'
import SignUp from './pages/SignUp'
import HomeList from './pages/homes/HomeList'
import HomeDetails from './pages/homes/HomeDetails'
import RoomDetails from './pages/rooms/RoomDetails'

const queryClient = new QueryClient()

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<SignUp />} />
            <Route path="/" element={<Layout children={undefined} />}>
              <Route index element={<Dashboard />} />
              <Route path="homes" element={<HomeList />} />
              <Route path="homes/:homeId" element={<HomeDetails />} />
              <Route path="homes/:homeId/rooms/:roomId" element={<RoomDetails />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App