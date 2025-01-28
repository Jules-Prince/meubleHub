import { 
    Box,
    Paper,
    Typography,
    TextField,
    Button,
    Container,
    Alert,
  } from '@mui/material';
  import { useState } from 'react';
  import { useNavigate } from 'react-router-dom';
  import { userService } from '../services/userService';
  
  interface SignUpFormData {
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
  }
  
  export default function SignUp() {
    const navigate = useNavigate();
    const [formData, setFormData] = useState<SignUpFormData>({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
    });
    const [error, setError] = useState<string>('');
    const [isLoading, setIsLoading] = useState(false);
  
    const handleSubmit = async (e: React.FormEvent) => {
      e.preventDefault();
      setError('');
  
      if (formData.password !== formData.confirmPassword) {
        setError('Passwords do not match');
        return;
      }
  
      setIsLoading(true);
  
      try {
        // We'll need to implement this in userService
        await userService.register({
          username: formData.username,
          email: formData.email,
          password: formData.password,
        });
        navigate('/login');
      } catch (error) {
        setError(error instanceof Error ? error.message : 'An error occurred');
      } finally {
        setIsLoading(false);
      }
    };
  
    return (
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          backgroundColor: '#f5f5f5'
        }}
      >
        <Container maxWidth="xs">
          <Paper 
            elevation={1} 
            sx={{ 
              p: 4,
              display: 'flex',
              flexDirection: 'column',
              gap: 2.5,
              borderRadius: 1,
            }}
          >
            <Typography 
              variant="h5" 
              component="h1"
              align="center"
              sx={{ 
                fontWeight: 500,
                mb: 1
              }}
            >
              MeubleHub
            </Typography>
  
            <Typography 
              variant="body1"
              align="center"
              color="text.secondary"
              sx={{ mb: 1 }}
            >
              Create your account
            </Typography>
  
            {error && (
              <Alert severity="error" sx={{ width: '100%' }}>
                {error}
              </Alert>
            )}
  
            <Box 
              component="form" 
              onSubmit={handleSubmit} 
              sx={{ 
                display: 'flex',
                flexDirection: 'column',
                gap: 2
              }}
            >
              <Box>
                <Typography
                  variant="caption"
                  component="label"
                  htmlFor="username"
                  sx={{
                    mb: 0.5,
                    display: 'block',
                    color: 'text.secondary'
                  }}
                >
                  Username
                </Typography>
                <TextField
                  id="username"
                  fullWidth
                  size="small"
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      bgcolor: 'white'
                    }
                  }}
                />
              </Box>
  
              <Box>
                <Typography
                  variant="caption"
                  component="label"
                  htmlFor="email"
                  sx={{
                    mb: 0.5,
                    display: 'block',
                    color: 'text.secondary'
                  }}
                >
                  Email Address
                </Typography>
                <TextField
                  id="email"
                  type="email"
                  fullWidth
                  size="small"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      bgcolor: 'white'
                    }
                  }}
                />
              </Box>
  
              <Box>
                <Typography
                  variant="caption"
                  component="label"
                  htmlFor="password"
                  sx={{
                    mb: 0.5,
                    display: 'block',
                    color: 'text.secondary'
                  }}
                >
                  Password
                </Typography>
                <TextField
                  id="password"
                  type="password"
                  fullWidth
                  size="small"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      bgcolor: 'white'
                    }
                  }}
                />
              </Box>
  
              <Box>
                <Typography
                  variant="caption"
                  component="label"
                  htmlFor="confirmPassword"
                  sx={{
                    mb: 0.5,
                    display: 'block',
                    color: 'text.secondary'
                  }}
                >
                  Confirm Password
                </Typography>
                <TextField
                  id="confirmPassword"
                  type="password"
                  fullWidth
                  size="small"
                  value={formData.confirmPassword}
                  onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      bgcolor: 'white'
                    }
                  }}
                />
              </Box>
  
              <Button
                type="submit"
                fullWidth
                variant="contained"
                disabled={isLoading}
                sx={{
                  textTransform: 'none',
                  bgcolor: '#1976d2',
                  '&:hover': {
                    bgcolor: '#1565c0'
                  }
                }}
              >
                Create Account
              </Button>
            </Box>
  
            <Box sx={{ 
              display: 'flex', 
              alignItems: 'center',
              gap: 2,
              my: 1
            }}>
              <Box sx={{ flex: 1, height: '1px', bgcolor: '#e0e0e0' }} />
              <Typography 
                variant="body2" 
                color="text.secondary"
              >
                OR
              </Typography>
              <Box sx={{ flex: 1, height: '1px', bgcolor: '#e0e0e0' }} />
            </Box>
  
            <Button
              fullWidth
              variant="outlined"
              onClick={() => navigate('/login')}
              sx={{
                textTransform: 'none',
                borderColor: '#e0e0e0',
                color: '#1976d2',
                '&:hover': {
                  borderColor: '#1976d2',
                  bgcolor: 'transparent'
                }
              }}
            >
              Sign in to existing account
            </Button>
          </Paper>
        </Container>
      </Box>
    );
  }