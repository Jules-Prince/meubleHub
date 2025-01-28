import { 
    AppBar,
    Toolbar,
    Typography,
    Button,
    IconButton,
    Box,
    Avatar,
    Menu,
    MenuItem,
    Divider,
  } from '@mui/material';
  import { useState, useEffect } from 'react';
  import { useNavigate } from 'react-router-dom';
  import LogoutIcon from '@mui/icons-material/Logout';
  import AccountCircleIcon from '@mui/icons-material/AccountCircle';
  import EmailIcon from '@mui/icons-material/Email';
  
  interface User {
    id: number;
    username: string;
    email: string;
  }
  
  export default function Navbar() {
    const navigate = useNavigate();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [user, setUser] = useState<User | null>(null);
    
    useEffect(() => {
      // Check localStorage for user data on component mount
      const userData = localStorage.getItem('user');
      if (userData) {
        setUser(JSON.parse(userData));
      }
    }, []);
  
    const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
      setAnchorEl(event.currentTarget);
    };
  
    const handleClose = () => {
      setAnchorEl(null);
    };
  
    const handleLogout = () => {
      localStorage.removeItem('user');
      setUser(null);
      handleClose();
      navigate('/login');
    };
  
    return (
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          {user ? (
            <>
              <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                MeubleHub
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <IconButton
                  size="large"
                  onClick={handleMenu}
                  color="inherit"
                >
                  <Avatar sx={{ width: 32, height: 32, bgcolor: 'secondary.main' }}>
                    {user.username[0].toUpperCase()}
                  </Avatar>
                </IconButton>
                <Menu
                  id="menu-appbar"
                  anchorEl={anchorEl}
                  anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  open={Boolean(anchorEl)}
                  onClose={handleClose}
                >
                  <Box sx={{ px: 2, py: 1 }}>
                    <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                      {user.username}
                    </Typography>
                    <Typography 
                      variant="body2" 
                      sx={{ 
                        display: 'flex', 
                        alignItems: 'center',
                        gap: 0.5,
                        color: 'text.secondary' 
                      }}
                    >
                      <EmailIcon fontSize="small" />
                      {user.email}
                    </Typography>
                  </Box>
                  <Divider />
                  <MenuItem onClick={handleLogout}>
                    <LogoutIcon sx={{ mr: 1 }} fontSize="small" />
                    Logout
                  </MenuItem>
                </Menu>
              </Box>
            </>
          ) : (
            <>
              <Box sx={{ flexGrow: 1, display: 'flex', alignItems: 'center' }}>
                <Typography variant="h6" component="div">
                  MeubleHub
                </Typography>
              </Box>
              <Button 
                color="inherit"
                onClick={() => navigate('/login')}
              >
                Login
              </Button>
            </>
          )}
        </Toolbar>
      </AppBar>
    );
  }