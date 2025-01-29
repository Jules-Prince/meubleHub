import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import Main from './pages/Main';
import Room from './pages/Room';
import Object from './pages/Object';
import LoginPage from './pages/Login';
import SignUpPage from './pages/Signin';
import HomeRooms from './pages/Room';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route
          path="/*"
          element={
            <Layout>
              <Routes>
                <Route path="/" element={<Main />} />
                <Route path="/home" element={<Home />} />
                <Route path="/homes/:id/rooms" element={<HomeRooms />} />
                <Route path="/room" element={<Room />} />
                <Route path="/object" element={<Object />} />
              </Routes>
            </Layout>
          }
        />
      </Routes>
    </Router>
  );
}

export default App;