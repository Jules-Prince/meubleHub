import { useState } from 'react';
import { 
  Box, 
  Typography, 
  Grid, 
  Button, 
  Card, 
  CardContent, 
  CardActions,
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  DialogActions
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import AddIcon from '@mui/icons-material/Add';
import HomeIcon from '@mui/icons-material/Home';

interface Home {
  id: number;
  name: string;
}

export default function HomeList() {
  const navigate = useNavigate();
  const [homes, setHomes] = useState<Home[]>([]);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [newHomeName, setNewHomeName] = useState('');

  const handleAddHome = () => {
    // TODO: API call to create home
    setIsAddDialogOpen(false);
    setNewHomeName('');
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">My Homes</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setIsAddDialogOpen(true)}
        >
          Add Home
        </Button>
      </Box>

      <Grid container spacing={3}>
        {homes.length === 0 ? (
          <Grid item xs={12}>
            <Typography align="center" color="text.secondary">
              You haven't created any homes yet. Click "Add Home" to get started.
            </Typography>
          </Grid>
        ) : (
          homes.map((home) => (
            <Grid item xs={12} sm={6} md={4} key={home.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
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
            </Grid>
          ))
        )}
      </Grid>

      <Dialog open={isAddDialogOpen} onClose={() => setIsAddDialogOpen(false)}>
        <DialogTitle>Add New Home</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Home Name"
            fullWidth
            value={newHomeName}
            onChange={(e) => setNewHomeName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setIsAddDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleAddHome} variant="contained">
            Create
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}