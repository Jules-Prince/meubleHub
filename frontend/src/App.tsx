import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import Main from './pages/Main';
import Room from './pages/Room';
import Object from './pages/Object';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Main />} />
          <Route path="/home" element={<Home />} />
          <Route path="/room" element={<Room />} />
          <Route path="/object" element={<Object />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;