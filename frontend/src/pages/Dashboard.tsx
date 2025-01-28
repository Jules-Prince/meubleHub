import { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper,
  Card,
  CardContent,
  CardActions,
  Button,
  Container
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import HomeIcon from '@mui/icons-material/Home';
import MeetingRoomIcon from '@mui/icons-material/MeetingRoom';
import ChairIcon from '@mui/icons-material/Chair';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';

interface DashboardStats {
  totalHomes: number;
  totalRooms: number;
  totalObjects: number;
  recentHomes: Array<{
    id: number;
    name: string;
  }>;
}

export default function Dashboard() {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    totalHomes: 0,
    totalRooms: 0,
    totalObjects: 0,
    recentHomes: []
  });

  // TODO: Fetch dashboard data when component mounts
  useEffect(() => {
    // API call will go here
  }, []);

  const StatCard = ({ icon, title, value }: { icon: React.ReactNode; title: string; value: number }) => (
    <Paper sx={{ p: 2, flex: 1 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
        {icon}
        <Typography variant="h6" sx={{ ml: 1 }}>
          {title}
        </Typography>
      </Box>
      <Typography variant="h4">{value}</Typography>
    </Paper>
  );

  return (
    <Container maxWidth="xl">
      <Typography variant="h4" gutterBottom>
        Welcome back, User
      </Typography>

      {/* Stats Overview */}
      <Box 
        sx={{ 
          display: 'flex', 
          gap: 3, 
          mb: 4,
          flexDirection: { xs: 'column', md: 'row' }
        }}
      >
        <StatCard 
          icon={<HomeIcon color="primary" />} 
          title="Total Homes" 
          value={stats.totalHomes} 
        />
        <StatCard 
          icon={<MeetingRoomIcon color="primary" />} 
          title="Total Rooms" 
          value={stats.totalRooms} 
        />
        <StatCard 
          icon={<ChairIcon color="primary" />} 
          title="Total Objects" 
          value={stats.totalObjects} 
        />
      </Box>

      {/* Recent Homes */}
      <Box sx={{ mb: 4 }}>
        <Box sx={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center', 
          mb: 2 
        }}>
          <Typography variant="h5">Recent Homes</Typography>
          <Button 
            endIcon={<ArrowForwardIcon />}
            onClick={() => navigate('/homes')}
          >
            View All Homes
          </Button>
        </Box>
        
        {stats.recentHomes.length === 0 ? (
          <Box sx={{ textAlign: 'center', py: 3 }}>
            <Typography color="text.secondary">
              You haven't created any homes yet. 
              <Button 
                color="primary" 
                onClick={() => navigate('/homes')}
                sx={{ ml: 1 }}
              >
                Create your first home
              </Button>
            </Typography>
          </Box>
        ) : (
          <Box sx={{ 
            display: 'grid',
            gap: 3,
            gridTemplateColumns: {
              xs: '1fr',
              sm: 'repeat(2, 1fr)',
              md: 'repeat(3, 1fr)',
              lg: 'repeat(4, 1fr)'
            }
          }}>
            {stats.recentHomes.map((home) => (
              <Card key={home.id}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <HomeIcon sx={{ mr: 1 }} />
                    <Typography variant="h6">{home.name}</Typography>
                  </Box>
                </CardContent>
                <CardActions>
                  <Button 
                    size="small" 
                    onClick={() => navigate(`/homes/${home.id}`)}
                  >
                    View Rooms
                  </Button>
                </CardActions>
              </Card>
            ))}
          </Box>
        )}
      </Box>

      {/* Quick Actions */}
      <Box>
        <Typography variant="h5" gutterBottom>
          Quick Actions
        </Typography>
        <Box sx={{ 
          display: 'flex', 
          gap: 2,
          flexDirection: { xs: 'column', sm: 'row' }
        }}>
          <Button 
            fullWidth 
            variant="contained" 
            onClick={() => navigate('/homes')}
            startIcon={<HomeIcon />}
          >
            Add New Home
          </Button>
          {stats.totalHomes > 0 && (
            <Button 
              fullWidth 
              variant="outlined"
              onClick={() => navigate(`/homes/${stats.recentHomes[0]?.id}`)}
              startIcon={<MeetingRoomIcon />}
            >
              Add New Room
            </Button>
          )}
        </Box>
      </Box>
    </Container>
  );
}