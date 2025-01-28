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
  DialogActions,
  Breadcrumbs,
  Link,
  Chip
} from '@mui/material';
import { useParams, Link as RouterLink } from 'react-router-dom';
import AddIcon from '@mui/icons-material/Add';
import ChairIcon from '@mui/icons-material/Chair';

interface Object {
  id: string;
  name: string;
  type: string;
  isReserved: boolean;
  reservedBy?: string;
}

export default function RoomDetails() {
  const { homeId, roomId } = useParams();
  const [objects, setObjects] = useState<Object[]>([]);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [newObject, setNewObject] = useState({
    name: '',
    type: ''
  });

  const handleAddObject = () => {
    // TODO: API call to create object
    setIsAddDialogOpen(false);
    setNewObject({ name: '', type: '' });
  };

  return (
    <Box>
      <Breadcrumbs sx={{ mb: 3 }}>
        <Link component={RouterLink} to="/homes" underline="hover">
          Homes
        </Link>
        <Link component={RouterLink} to={`/homes/${homeId}`} underline="hover">
          Current Home
        </Link>
        <Typography color="text.primary">Current Room</Typography>
      </Breadcrumbs>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">Objects</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setIsAddDialogOpen(true)}
        >
          Add Object
        </Button>
      </Box>

      <Grid container spacing={3}>
        {objects.length === 0 ? (
          <Grid item xs={12}>
            <Typography align="center" color="text.secondary">
              No objects added yet. Click "Add Object" to create one.
            </Typography>
          </Grid>
        ) : (
          objects.map((object) => (
            <Grid item xs={12} sm={6} md={4} key={object.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <ChairIcon sx={{ mr: 1 }} />
                    <Typography variant="h6">{object.name}</Typography>
                  </Box>
                  <Typography color="text.secondary" gutterBottom>
                    Type: {object.type}
                  </Typography>
                  {object.isReserved && (
                    <Chip 
                      label={`Reserved by ${object.reservedBy}`}
                      color="primary"
                      size="small"
                    />
                  )}
                </CardContent>
                <CardActions>
                  {!object.isReserved && (
                    <Button size="small" color="primary">
                      Reserve
                    </Button>
                  )}
                </CardActions>
              </Card>
            </Grid>
          ))
        )}
      </Grid>

      <Dialog open={isAddDialogOpen} onClose={() => setIsAddDialogOpen(false)}>
        <DialogTitle>Add New Object</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Object Name"
            fullWidth
            value={newObject.name}
            onChange={(e) => setNewObject({ ...newObject, name: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Object Type"
            fullWidth
            value={newObject.type}
            onChange={(e) => setNewObject({ ...newObject, type: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setIsAddDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleAddObject} variant="contained">
            Create
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}