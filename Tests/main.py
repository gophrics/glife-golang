import ProfileService
import TravelService
import Context

def main():
    ProfileService.username_exist()
    success = ProfileService.register()
    if success:
        ProfileService.get_me()
        ProfileService.get_user()
        TravelService.get_all_trips('oynlphjur')
        TravelService.save_trip()

if __name__ == "__main__":
    main()
